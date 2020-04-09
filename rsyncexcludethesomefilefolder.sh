[root@ip-10-93-32-34 TAP86]# cat /opt/package.sh
#!/bin/bash

echo "Running RSYNC"

rsync -a --exclude={'.bash_history','.bash_logout','.bash_profile','.bashrc','.ssh','.viminfo',} /home/tmsi_admin/ /var/www/nginx-default/tmsi/TAP86/

echo "Changing permission"

chown -R nginx:nginx /var/www/nginx-default/tmsi/TAP86/