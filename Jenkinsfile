def registryHost = ""
def tagName = ""
def targetPath = ""
pipeline {
    agent any
    stages  {
        stage("检查构建分支") {
            steps {
                echo "检查构建分支中......"
                script {
                    APP_NAME = "$APP_NAME"
                    if (APP_NAME ==~ /^api-.*/) {
                       targetPath = "$WORKSPACE/$APP_NAME"
                    } else {
                        targetPath = "$WORKSPACE/service/$APP_NAME"
                    }
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/)  {
                        echo "构建正式环境，tag=${env.GIT_BRANCH}"
                        tagName = env.GIT_BRANCH
                        registryHost = env.PRO_REGISTRY_HOST
                    } else if (env.GIT_BRANCH ==~ /^release-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        echo "构建预生产环境，tag=${env.GIT_BRANCH}"
                        tagName = env.GIT_BRANCH
                    } else if (env.GIT_BRANCH ==~ /(origin\/develop)/) {
                        registryHost = env.DEV_REGISTRY_HOST
                        tagName="develop"
                        echo "构建开发环境，/origin/develop"
                    } else {
                        echo "构建分支${env.GIT_BRANCH}不合法，只允许构建正式环境分支(例如：v1.0.0)，预生产环境分支(例如：release-1.0.0)，开发环境分支(/origin/develop)"
                        sh 'exit 1'
                    }
                }
                echo "检查构建分支完成."
            }
        }

        stage("代码检查") {
            steps {
                echo "代码检查中......"
                //sh "$goLintcmd"
                echo "代码检查完成."
            }
        }

        stage("代码编译") {
            steps {
                echo "代码编译中......"
                script {
                     sh "cd $targetPath && bash ./scripts/build-app.sh $gocmd"
                }
                echo "代码编译完成."
            }
        }

        stage("单元测试") {
            steps {
                echo "单元测试中......"
                echo "单元测试完成."
            }
        }

        stage("集成测试") {
            steps {
                echo "集成测试......"
                echo "集成测试完成"
            }
        }

        stage("构建镜像") {
            steps {
                echo "构建镜像中......"
                // 兼容自动构建和参数构建
                script {
                    sh "cd $targetPath && bash ./scripts/build-image.sh $registryHost $tagName $APP_NAME $configPath"
                }
                echo "构建镜像完成"
            }
        }

        stage("上传镜像") {
            steps {
                echo "上传镜像中......"
                // 兼容自动构建和参数构建
                script {
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        echo "使用正式环境镜像仓库 ${registryHost}"
                    }
                    else if (env.GIT_BRANCH ==~ /^release-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                          echo "使用预生产环境镜像仓库 ${registryHost}"
                    }
                    else {
                        echo "使用开发环境 ${registryHost}"
                    }
                    sh "pwd"
                    sh "cd $targetPath && bash ./scripts/push-image.sh $registryHost $tagName $APP_NAME"
                }
                echo "上传镜像完成"
            }
        }

        stage("部署到远程服务器") {
            // 正式环境和预生产环境跳过部署，手动部署
//             when { expression { return env.GIT_BRANCH ==~ /(origin\/staging|origin\/develop)/ } }
            steps {
                echo "部署到远程服务器"
                script {
                   sh "pwd"
                   docker_image="$APP_NAME-$tagName"
                   server_ip=env.SERVER_HOST
                   log_host=env.LOG_SERVER_HOST
                   sh "cd $targetPath && bash ./scripts/deploy.sh $registryHost $docker_image $server_ip $APP_NAME"
                }
                echo "部署到远程服务器完成"
            }
        }

        stage("清理空间") {
            steps {
                echo "清理空间......"
                   script {
                      sh "pwd"
                      try {
                         sh "rm go.mod"
                      }catch(Exception e) {
                          println e
                      }
                      //sh "cd $target && bash ./scripts/delete-images.sh $APP_NAME"
                   }
                echo "清理空间完成"
            }
        }
    }

    post {
    		always {
    			echo 'One way or another, I have finished'
    			echo sh(returnStdout: true, script: 'env')
    			deleteDir() /* clean up our workspace */
    		}
    		success {
    			SendDingding("success")
    			echo 'structure success'
    		}
    		failure {
    			SendDingding("failure")
    			echo 'structure failure'
    		}
    		unsuccessful {
    			//SendDingding("unsuccessful")
    			echo 'structure unsuccessful'
    		}
    		aborted {
    			//SendDingding("aborted")
    			echo 'structure aborted'
    		}
    		unstable {
    			//SendDingding("unstable")
    			echo 'structure unstable'
    		}
       }
}

void SendEmail(res)
{
	//在这里定义邮箱地址
	addr="xxx@xxx.com"
	if( res == "success" )
	{
		mail to: addr,
		subject: "构建成功 ：${currentBuild.fullDisplayName}",
		body: "\n项目顺利构建成功，恭喜你！ \n\n任务名称： ${env.JOB_NAME} 第 ${env.BUILD_NUMBER} 次构建 \n\n 更多信息请查看 : ${env.BUILD_URL}"
	}
	else
	{
		mail to: addr,
		subject: "构建失败 ：${currentBuild.fullDisplayName}",
		body: "\n完蛋了，快去看一下出了啥问题！ \n\n任务名称： ${env.JOB_NAME} 第 ${env.BUILD_NUMBER} 次构建 \n\n 更多信息请查看 : ${env.BUILD_URL}"
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
