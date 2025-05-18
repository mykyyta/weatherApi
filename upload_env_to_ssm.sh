#!/bin/bash

# Prefix for all parameters in SSM
PARAM_PREFIX="/weather-api"

# Path to .env file
ENV_FILE=".env"

# AWS settings
AWS_PROFILE=${AWS_PROFILE:-default}
AWS_REGION=${AWS_REGION:-us-east-1}

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
  echo "❌ $ENV_FILE not found!"
  exit 1
fi

echo "🔁 Uploading environment variables from $ENV_FILE"
echo "📍 Using profile: $AWS_PROFILE | region: $AWS_REGION"
echo "📦 Parameter prefix: $PARAM_PREFIX"
echo

# Read .env file line by line
while IFS='=' read -r key value || [[ -n "$key" ]]; do
  # Skip comments and empty lines
  if [[ "$key" =~ ^#.*$ || -z "$key" ]]; then
    continue
  fi

  # Trim whitespace
  key=$(echo "$key" | xargs)
  value=$(echo "$value" | xargs)

  # Validate key/value
  if [ -z "$key" ] || [ -z "$value" ]; then
    echo "⚠️ Skipping invalid line: $key=$value"
    continue
  fi

  PARAM_NAME="${PARAM_PREFIX}/${key}"

  echo "➕ Uploading: $PARAM_NAME"
  aws ssm put-parameter \
    --name "$PARAM_NAME" \
    --value "$value" \
    --type "SecureString" \
    --overwrite \
    --region "$AWS_REGION" \
    --profile "$AWS_PROFILE"
done < "$ENV_FILE"

echo
echo "✅ Done uploading environment variables to SSM Parameter Store."