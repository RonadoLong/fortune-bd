pipeline{
    agent any
    environment{
        HARBOR_HOST='192.168.3.30:8086'
        HARBOR_ADDR='192.168.3.30:8086/mateforce'
        K8S_NAMESPACE='develop'
    }
    parameters {
//         string(name: 'PROJECT_NAME', defaultValue: '', description: 'project name,same as the name ofdocker container')
//         string(name: 'CONTAINER_VERSION', defaultValue: '', description: 'docker container version number, SET when major version number changed')
        booleanParam(name: 'DEPLOYMENT_K8S', defaultValue: false, description: 'release deployment k8s')
    }
    stages {
        stage('Initial') {
            steps{
                script {
                        env.DOCKER_IMAGE='${PROJECT_NAME}'
                         APP_NAME = "$PROJECT_NAME"
                         if (APP_NAME ==~ /^api-.*/) {
                             env.TARGET_PATH = "./${APP_NAME}"
                         } else {
                            env.TARGET_PATH = "./app/${APP_NAME}"
                         }
                          // 脚本式创建一个环境变量
                        if (params.CONTAINER_VERSION == '') {
                                env.APP_VERSION = 'v1.0.0-alpha'
//                             env.APP_VERSION = sh(returnStdout:true,script:"jenkins-build-tools gen -p ${params.PROJECT_NAME}").trim()
                        }else {
                            env.APP_VERSION ="${params.CONTAINER_VERSION}-alpha"
                        }
                        sh "echo ${env.APP_VERSION}"
                    }
                }
        }
        stage("Docker Build") {
            when {
                allOf {
                    expression { env.APP_VERSION != null }
                }
            }
            steps("Start Build") {
//                 sh "docker login -u admin -p QQabc123++ ${HARBOR_HOST}"
                sh "docker build --build-arg TARGET_PATH=${TARGET_PATH} -t ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION} -f ${TARGET_PATH}/deploy/Dockerfile ."
                sh "docker push ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION}"
                sh "docker rmi ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION} -f"
            }

        }
        stage("Deploy") {
            when {
                allOf {
                    expression { env.APP_VERSION != null }
                }
            }
            steps("Deploy to kubernetes") {
                script {
                    if (params.DEPLOYMENT_K8S) {
                        sh "export KUBECONFIG=${env.KUBECONFIG}"
                        sh "sed -i 's/VERSION_NUMBER/${APP_VERSION}/g' ${TARGET_PATH}/deploy/k8s-deployment.yml"
                        sh "kubectl apply -f ${TARGET_PATH}/deploy/k8s-deployment.yml --namespace=develop"
                    }
                }
            }
        }
    }
    post {
    		always {
    			echo 'One way or another, I have finished'
//     			echo sh(returnStdout: true, script: 'env')
    			deleteDir() /* clean up our workspace */
    		}
    		success {
//     			SendDingding("success")
    			echo 'structure success'
    		}
    		failure {
//     			SendDingding("failure")
    			echo 'structure failure'
    		}
       }
}

void SendDingding(res)
{
	// 输入相应的手机号码，在钉钉群指定通知某个人
	tel_num="13377000902"

	// 钉钉机器人的地址
	dingding_url="https://oapi.dingtalk.com/robot/send\\?access_token\\=a5e1c5e003cc109af1ad0ff85a0d4ca35fee804da0b7e77c156fd24c9bc6a16d"

    branchName=""
    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
        branchName="理财项目正式环境 tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else if (env.GIT_BRANCH ==~ /^release-([0-9])+\.([0-9])+\.([0-9])+.*/){
        branchName="理财项目预生产环境 tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else {
        branchName="理财项目开发环境 branch=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }

    // 发送内容
	json_msg=""
	if( res == "success" ) {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [送花花] ' + "${branchName} 第${env.BUILD_NUMBER}次构建，"  + '构建成功。 \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}
	else {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [大哭] ' + "${branchName} 第${env.BUILD_NUMBER}次构建，"  + '构建失败，请及时处理！ \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}

    post_header="Content-Type:application/json;charset=utf-8"
    sh_cmd="curl -X POST " + dingding_url + " -H " + "\'" + post_header + "\'" + " -d " + "\""  + json_msg + "\""
// 	sh sh_cmd
}
