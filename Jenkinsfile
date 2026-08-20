pipeline {
    agent {
        label 'ubuntu'
    }

    triggers {
        pollSCM('H/5 * * * *')
    }

    options {
        timeout(time: 1, unit: 'HOURS')
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '30', artifactNumToKeepStr: '10'))
        disableConcurrentBuilds()
    }

    environment {
        REGISTRY              = 'docker.io'
        IMAGE_NAME            = 'mycompany/order-service'
        DOCKER_CREDENTIALS_ID = 'docker-registry-credentials'
        GIT_CREDENTIALS_ID    = 'git-credentials'
        COVERAGE_THRESHOLD    = '70'
    }

    stages {
        stage('Checkout & Environment Init') {
            steps {
                script {
                    echo "================================================="
                    echo "STAGE 1: Checkout & Environment Initialization"
                    echo "================================================="

                    def rawVersion = readFile('VERSION').trim()
                    def commitSha  = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()

                    env.VERSION_TAG      = "${rawVersion}.${BUILD_NUMBER}"
                    env.COMMIT_SHA       = commitSha
                    env.FULL_IMAGE_TAG   = "${REGISTRY}/${IMAGE_NAME}:${env.VERSION_TAG}"
                    env.LATEST_IMAGE_TAG = "${REGISTRY}/${IMAGE_NAME}:latest"
                    env.CURRENT_BRANCH   = env.BRANCH_NAME ?: env.GIT_BRANCH ?: 'main'

                    echo "Build Branch      : ${env.CURRENT_BRANCH}"
                    echo "Build Version Tag : ${env.VERSION_TAG}"
                    echo "Git Commit SHA    : ${env.COMMIT_SHA}"
                    echo "Target Image Tag  : ${env.FULL_IMAGE_TAG}"
                }
            }
        }

        stage('Static Analysis & Linting') {
            steps {
                echo "================================================="
                echo "STAGE 2: Code Quality & Static Analysis (PR & Main)"
                echo "================================================="
                sh 'golangci-lint run --timeout 5m ./...'
            }
        }

        stage('SAST Security Scan') {
            steps {
                echo "================================================="
                echo "STAGE 3: SAST Security & Vulnerability Audit (PR & Main)"
                echo "================================================="
                sh '''
                    set +e
                    command -v gosec >/dev/null 2>&1 || go install github.com/securego/gosec/v2/cmd/gosec@latest
                    command -v govulncheck >/dev/null 2>&1 || go install golang.org/x/vuln/cmd/govulncheck@latest

                    echo "--> Running Gosec security scanner..."
                    gosec -fmt=text ./... || true

                    echo "--> Running Go Vulnerability Auditor (govulncheck)..."
                    govulncheck ./... || true
                    exit 0
                '''
            }
        }

        stage('Unit Tests & Code Coverage') {
            steps {
                echo "================================================="
                echo "STAGE 4: Unit Testing & Code Coverage Check (PR & Main)"
                echo "================================================="
                sh '''
                    go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

                    COVERAGE=$(go tool cover -func=coverage.out | grep total: | awk '{print $3}' | sed 's/%//')
                    echo "Total Code Coverage: ${COVERAGE}%"

                    if [ $(echo "${COVERAGE} < ${COVERAGE_THRESHOLD}" | bc -l 2>/dev/null || awk "BEGIN {print (${COVERAGE} < ${COVERAGE_THRESHOLD})}") -eq 1 ]; then
                        echo "[WARNING] Coverage ${COVERAGE}% is below target threshold ${COVERAGE_THRESHOLD}%"
                    fi
                '''
            }
            post {
                always {
                    archiveArtifacts artifacts: 'coverage.out', allowEmptyArchive: true
                }
            }
        }

        stage('Application Compile Check') {
            steps {
                echo "================================================="
                echo "STAGE 5: Compiling Binary Verification (PR & Main)"
                echo "================================================="
                sh '''
                    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
                        -ldflags="-s -w -X 'main.Version=${VERSION_TAG}'" \
                        -trimpath \
                        -o bin/server ./cmd/server
                '''
            }
        }

        stage('Docker Build') {
            when {
                anyOf {
                    branch 'main'
                    expression { env.BRANCH_NAME == 'main' }
                    expression { env.GIT_BRANCH == 'main' }
                    expression { env.GIT_BRANCH == 'origin/main' }
                }
            }
            steps {
                echo "================================================="
                echo "STAGE 6: Building Docker Container Image (Main Branch Only)"
                echo "================================================="
                sh '''
                    docker build \
                        --build-arg VERSION="${VERSION_TAG}" \
                        --build-arg COMMIT_SHA="${COMMIT_SHA}" \
                        -t "${FULL_IMAGE_TAG}" \
                        -t "${LATEST_IMAGE_TAG}" .
                '''
            }
        }

        stage('Container Security Scan (Trivy)') {
            when {
                anyOf {
                    branch 'main'
                    expression { env.BRANCH_NAME == 'main' }
                    expression { env.GIT_BRANCH == 'main' }
                    expression { env.GIT_BRANCH == 'origin/main' }
                }
            }
            steps {
                echo "================================================="
                echo "STAGE 7: Container Vulnerability Scan (Main Branch Only)"
                echo "================================================="
                sh '''
                    trivy image --exit-code 0 --severity HIGH,CRITICAL --format table ${FULL_IMAGE_TAG}
                '''
            }
        }

        stage('Push to Docker Registry') {
            when {
                anyOf {
                    branch 'main'
                    expression { env.BRANCH_NAME == 'main' }
                    expression { env.GIT_BRANCH == 'main' }
                    expression { env.GIT_BRANCH == 'origin/main' }
                }
            }
            steps {
                echo "================================================="
                echo "STAGE 8: Pushing Container Image to Registry (Main Branch Only)"
                echo "================================================="
                withCredentials([usernamePassword(credentialsId: env.DOCKER_CREDENTIALS_ID, usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh '''
                        echo "${DOCKER_PASS}" | docker login ${REGISTRY} -u "${DOCKER_USER}" --password-stdin
                        docker push "${FULL_IMAGE_TAG}"
                        docker push "${LATEST_IMAGE_TAG}"
                        docker logout ${REGISTRY}
                    '''
                }
            }
        }

        stage('Update Git Image Tag') {
            when {
                anyOf {
                    branch 'main'
                    expression { env.BRANCH_NAME == 'main' }
                    expression { env.GIT_BRANCH == 'main' }
                    expression { env.GIT_BRANCH == 'origin/main' }
                }
            }
            steps {
                echo "================================================="
                echo "STAGE 9: Updating Git Manifest Image Tag (Main Branch Only)"
                echo "================================================="
                withCredentials([sshUserPrivateKey(credentialsId: env.GIT_CREDENTIALS_ID, keyFileVariable: 'SSH_KEY', usernameVariable: 'GIT_USER')]) {
                    sh '''
                        chmod +x ./scripts/update-image-tag.sh
                        ./scripts/update-image-tag.sh "${VERSION_TAG}" "deployments/k8s/deployment.yaml"
                    '''
                }
            }
        }
    }

    post {
        always {
            echo "Cleaning up workspace..."
            cleanWs()
        }
        success {
            echo "SUCCESS: Jenkins Pipeline successfully completed for branch ${env.CURRENT_BRANCH}"
        }
        failure {
            echo "FAILURE: Jenkins Pipeline failed at build #${BUILD_NUMBER}"
        }
    }
}
