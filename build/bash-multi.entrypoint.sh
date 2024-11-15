#!/usr/bin/env sh

#if [ ! -d "/log/be-auth-service" ]; then
  #mkdir -p /log/auth
#fi

#get aws s3 configs
aws configure set role_profile arn:aws:iam::176300518568:role/dev-jenkins-slave-admin-role
aws configure set source_profile default

if [ $ENVIRONMENT = "prd" ]; then
    aws s3 sync --no-progress s3://servicehub-be-configs/${ENVIRONMENT}/be-${SERVICE}/ conf/ --profile default
else
    aws s3 sync --no-progress s3://servicehub-be-configs/dev/${ENVIRONMENT}/be-${SERVICE}/ conf/ --profile default
fi



#aws s3 sync --no-progress s3://servicehub-be-configs/dev/${ENVIRONMENT}/be-${SERVICE}/ conf/ --profile default

# Execute the original command.
exec "$@"
