@ECHO OFF
SET DIRNAME=%~dp0
"java.exe" -Xmx64m -Xms64m -classpath "%DIRNAME%gradle\wrapper\gradle-wrapper.jar" org.gradle.wrapper.GradleWrapperMain %*
