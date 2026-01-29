#!/bin/env fish

if test -z "$SECRETS_PATH" || test -z "$ENC_SECRETS_PATH" || test -z "$AGE_KEY"
    exit 1
end

if test $argv[1] && test $argv[1] = enc
    for file in (ls $SECRETS_PATH)
        echo Encrypting: $file into $ENC_SECRETS_PATH/$file.age
        age --encrypt \
            --output $ENC_SECRETS_PATH/$file.age \
            --identity $AGE_KEY \
            $SECRETS_PATH/$file
    end
else if test $argv[1] && test $argv[1] = dec
    mkdir -p $SECRETS_PATH
    for file in (ls $ENC_SECRETS_PATH)
        set new_file (path change-extension '' $file)
        echo Decrypting $file into $SECRETS_PATH/$new_file
        age --decrypt \
            --identity $AGE_KEY \
            $ENC_SECRETS_PATH/$file >$SECRETS_PATH/$new_file
    end
else
    echo "need arg: enc or dec"
    exit 1
end
