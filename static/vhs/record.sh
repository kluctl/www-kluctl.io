#!/usr/bin/env bash

set -e

DIR=$(cd $(dirname $0) && pwd)
cd $DIR

if [ ! -z "$(git status --porcelain -- ./demo-project)" ]; then
  echo "demo-project must be clean"
  git status --porcelain -- ./demo-project
  exit 1
fi

ORIG_BRANCH=$(git rev-parse --abbrev-ref HEAD)
git checkout -B demo-branch
function cleanup() {
  # undo changes
  git checkout -- demo-project

  echo "restoring git"
  git checkout $ORIG_BRANCH

  echo webui_pid=$webui_pid
  if [ "$webui_pid" != "" ]; then
    kill $webui_pid
  fi
}
trap cleanup EXIT

git push origin demo-branch -f

if ! kind get clusters | grep "kluctl-demo-tape" > /dev/null; then
  kind create cluster -n kluctl-demo-tape
fi
export KLUCTL_CONTEXT=kind-kluctl-demo-tape

echo "Waiting for cluster..."
while ! kubectl get pod -A 2>&1 >/dev/null; do sleep 1; done

kluctl controller install --yes
echo "waiting for controller..."
kubectl -n kluctl-system wait --for=condition=available --timeout=1h deployments/kluctl-controller

# we don't want to see the controller installation
kubectl delete ns kluctl-results || true

# delete old run
kubectl delete ns test || true
kubectl -n kluctl-system delete kluctldeployment demo || true

kluctl webui run --port=9090 &
webui_pid=$!
sleep 5

rm -rf recordings
node record-webui.js

mv demo-project/demo-cli.mp4 .
mv recordings/*.webm ./demo-webui.webm

rm -f demo-webui.mp4
ffmpeg -i demo-webui.webm demo-webui.mp4
rm demo-webui.webm
