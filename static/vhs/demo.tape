Set Shell zsh

Set WindowBar Colorful

Set Width 1000
Set Height 800

Set Margin 10
Set MarginFill "#DDDDDD"
Set BorderRadius 15
Set Padding 15

Type "# Lets do a fresh deployment to test."
Ctrl+C
Sleep 1s

Type "kluctl deploy -t test --yes"
Enter
Sleep 3s

Enter
Type "# Lets check our deployment."
Ctrl+C
Type "kubectl -n test get pod"
Enter
Sleep 2s

Enter
Type "# Lets modify something (scale nginx)."
Ctrl+C

Type "yq -i '.spec.replicas=3' deployment/nginx/deploy.yml"
Enter
Sleep 1s

Enter
Type "# Lets see a diff between the actual and desired state."
Ctrl+C
Sleep 2s

Type "kluctl diff -t test"
Enter
Sleep 3s

Enter
Type "# Looks reasonable, lets deploy it!"
Ctrl+C

Type "kluctl deploy -t test --yes"
Enter
Sleep 3s

Enter
Type "# Lets check our deployment. You should see 3 pods now."
Ctrl+C
Type "kubectl -n test get pod"
Enter
Sleep 2s

Enter
Type "# It did what it promised to do :)"
Ctrl+C
Type "# Now let GitOps take over!"
Ctrl+C
Sleep 2s

Type "kubectl apply -f kluctl-gitops.yaml"
Enter
Sleep 0.5s
Type "# Waiting for the reconciliation in the background..."
Ctrl+C
Sleep 5s

Enter
Type "# Check our deployment again, also look at the Webui!"
Ctrl+C
Type "kubectl -n test get pod"
Enter
Sleep 2s

Enter
Type "# You should see that it has reverted the change we did before."
Ctrl+C
Type "# This is because we did not push it to Git."
Ctrl+C
Enter
Type "# Let's make some additional changes (a bugfix) and push it now!"
Ctrl+C
Sleep 1s

Set TypingSpeed 20ms
Type "yq -i '.spec.template.spec.containers[0].ports[0].containerPort=8080' deployment/nginx/deploy.yml"
Enter
Type "yq -i '.spec.ports[0].name="
Type '"http"'
Type "' deployment/nginx/service.yml"
Enter
Type "yq -i '.spec.ports[0].targetPort="
Type '"http"'
Type "' deployment/nginx/service.yml"
Enter
Type "yq -i '.spec.selector.app="
Type '"http"'
Type " | del(.spec.selector.appp)' deployment/nginx/service.yml"
Enter
Sleep 0.5s
Set TypingSpeed 50ms

Type "git add deployment/nginx/*.yml"
Enter
Sleep 0.5s
Type "git commit -m 'some changes'"
Enter
Sleep 0.5s
Type "git push origin demo-branch"
Enter
Sleep 4s
Type "# Waiting for the reconciliation in the background..."
Ctrl+C
Sleep 5s

Enter
Type "# Check our deployment in the Webui."
Ctrl+C
Type "# You should see our bugfix being applied."
Ctrl+C
Type "# You should also see the replicas being 3 again"
Ctrl+C
Type "kubectl -n test get pod"
Enter
Sleep 2s

Sleep 10s

Ctrl+D

Output demo-cli.mp4
