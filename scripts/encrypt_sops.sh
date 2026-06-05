#!/bin/sh

set -eu

if [ -z "${AGE_KEY_NO_PQ}" ]; then
	echo "unbound variable"
fi
if [ ! -f "${AGE_KEY_NO_PQ}" ]; then
	echo "Error: ${AGE_KEY_NO_PQ} file does not exist"
	exit 1
fi

PUBLIC_KEY=$(age-keygen -y $AGE_KEY_NO_PQ)

SECRETS_ENC_PATH=$KUBE_ENC_SECRETS
mkdir -p $SECRETS_ENC_PATH
for FILE in $KUBE_SECRETS/*; do
	FILENAME="${FILE##*/}"
	DEST=$SECRETS_ENC_PATH/$FILENAME

	if [ "$FILENAME" = "kustomization.yaml" ]; then
		mv $FILE $DEST
		echo "Moving UNENCRYPTED $FILE"
		continue
	fi

	sops --encrypt --in-place $FILE
  echo Moving encrypted file to $DEST
  mv $FILE $DEST
done
