#!/usr/bin/env bash
# ==============================================================================
# Enterprise Git Ops Image Tag Update Script
# Updates deployment manifest with new Docker image tag and commits back to Git.
# ==============================================================================

set -euo pipefail

NEW_TAG="${1:-}"
MANIFEST_FILE="${2:-deployments/k8s/deployment.yaml}"
GIT_USER_NAME="${GIT_USER_NAME:-Jenkins CI Bot}"
GIT_USER_EMAIL="${GIT_USER_EMAIL:-jenkins-ci@example.com}"

if [ -z "${NEW_TAG}" ]; then
  echo "[ERROR] Missing required argument: NEW_TAG"
  echo "Usage: $0 <NEW_TAG> [MANIFEST_FILE]"
  exit 1
fi

if [ ! -f "${MANIFEST_FILE}" ]; then
  echo "[ERROR] Deployment manifest file not found: ${MANIFEST_FILE}"
  exit 1
fi

echo "================================================="
echo "Updating Image Tag in Git Repository"
echo "Target Manifest : ${MANIFEST_FILE}"
echo "New Image Tag   : ${NEW_TAG}"
echo "================================================="

# Update image tag in Kubernetes deployment manifest using sed
# Matches 'image: <any_registry>/<image_name>:<old_tag>' and replaces tag
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' -E "s|(image: .*:)[^[:space:]]+|\1${NEW_TAG}|g" "${MANIFEST_FILE}"
else
  sed -i -E "s|(image: .*:)[^[:space:]]+|\1${NEW_TAG}|g" "${MANIFEST_FILE}"
fi

echo "[INFO] Manifest updated successfully."

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "[WARN] Not inside a Git repository. Skipping git commit and push."
  exit 0
fi

echo "[INFO] Verifying git diff:"
git diff "${MANIFEST_FILE}"

# Configure Git user identity for automated commit
git config user.name "${GIT_USER_NAME}"
git config user.email "${GIT_USER_EMAIL}"

# Check if there are changes to commit
if git diff --quiet "${MANIFEST_FILE}"; then
  echo "[WARN] No changes detected in ${MANIFEST_FILE}. Image tag is already set to ${NEW_TAG}."
  exit 0
fi

# Stage, Commit, and Push changes
echo "[INFO] Staging and committing manifest update..."
git add "${MANIFEST_FILE}"

COMMIT_MSG="[ci skip] chore(deploy): update image tag to ${NEW_TAG}"
git commit -m "${COMMIT_MSG}"

# Optional: Create Git tag for release traceability
RELEASE_TAG="v${NEW_TAG}"
if ! git rev-parse "${RELEASE_TAG}" >/dev/null 2>&1; then
  echo "[INFO] Creating Git Tag: ${RELEASE_TAG}"
  git tag -a "${RELEASE_TAG}" -m "Release ${RELEASE_TAG}"
fi

# Push changes to remote repository
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "[INFO] Pushing changes to branch '${CURRENT_BRANCH}' and tags..."

if [ "${DRY_RUN:-false}" = "true" ]; then
  echo "[DRY RUN] Would execute: git push origin ${CURRENT_BRANCH} --tags"
else
  git push origin "${CURRENT_BRANCH}" --tags || echo "[WARN] Git push failed. Verify git credentials & remote origin."
fi

echo "[SUCCESS] Image tag update completed successfully!"

