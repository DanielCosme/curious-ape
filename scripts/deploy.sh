#!/usr/bin/env fish

if test "$RELEASE" = true
  echo "--- Pushing code to remote ---"
  set current_branch (git branch --show-current)

  if not test $current_branch = "main"
    echo "git is NOT on the main branch"
    echo "current branch: " $current_branch
    exit 1
  end

  git diff --exit-code; or echo "Working tree cannot be dirty" and exit 1
  git push; or exit 1
  echo "--- Success ---"
  echo ""

  echo "--- Starting Release ---"
  echo "\
    cd curious-ape
    git checkout main; or exit 1
    git pull; or exit 1
    ./scripts/release.sh; or exit 1 \
    " | ssh daniel@ubi-prime ; or exit 1
    echo "--- Success ---"
  echo ""
end


# We are assuming that the kubectl client is properly configured.
echo "--- Synchronizing Kubernetes resources ---"

kubectl create \
        configmap \
        curious-ape-prod \
        --from-file=./kube/prod/config.json \
        --dry-run=client \
        -o yaml \
        | kubectl apply -f -

kubectl apply -f ./kube/deployment.yaml
kubectl rollout restart deployment curious-ape

echo "--- Success ---"