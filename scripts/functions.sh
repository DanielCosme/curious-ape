do_ssh() {
  ssh $USER@$SERVER $SSH_OPTIONS $@
}

do_scp() {
  scp $SSH_OPTIONS $1 $USER@$SERVER:$2
}