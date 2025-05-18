#!/bin/bash

PARAM_PREFIX="/weather-api"
ENV_FILE=".env"
AWS_PROFILE=${AWS_PROFILE:-default}
AWS_REGION=${AWS_REGION:-us-east-1}

if [ ! -f "$ENV_FILE" ]; then
  echo "❌ $ENV_FILE not found!"
  exit 1
fi

echo "🔁 Uploading environment variables from $ENV_FILE"
echo "📍 Using profile: $AWS_PROFILE | region: $AWS_REGION"
echo "📦 Parameter prefix: $PARAM_PREFIX"
echo

while IFS='=' read -r key value || [[ -n "$key" ]]; do
  if [[ "$key" =~ ^#.*$ || -z "$key" ]]; then
    continue
  fi

  key=$(echo "$key" | xargs)
  value=$(echo "$value" | xargs)

  if [ -z "$key" ] || [ -z "$value" ]; then
    echo "⚠️ Skipping invalid line: $key=$value"
    continue
  fi

  PARAM_NAME="${PARAM_PREFIX}/${key}"

  echo "🗑️  Deleting existing: $PARAM_NAME (if any)"
  aws ssm delete-parameter \
    --name "$PARAM_NAME" \
    --region "$AWS_REGION" \
    --profile "$AWS_PROFILE" 2>/dev/null

  echo "➕ Uploading as version 1: $PARAM_NAME"
  aws ssm put-parameter \
    --name "$PARAM_NAME" \
    --value "$value" \
    --type "SecureString" \
    --region "$AWS_REGION" \
    --profile "$AWS_PROFILE"
done < "$ENV_FILE"

# Ensure GIN_MODE is set to release if not in .env
GIN_MODE_VALUE=$(grep '^GIN_MODE=' "$ENV_FILE" | cut -d '=' -f2-)
if [ -z "$GIN_MODE_VALUE" ]; then
  PARAM_NAME="${PARAM_PREFIX}/GIN_MODE"
  echo "🗑️  Deleting existing: $PARAM_NAME (if any)"
  aws ssm delete-parameter \
    --name "$PARAM_NAME" \
    --region "$AWS_REGION" \
    --profile "$AWS_PROFILE" 2>/dev/null

  echo "➕ Uploading: $PARAM_NAME = release"
  aws ssm put-parameter \
    --name "$PARAM_NAME" \
    --value "release" \
    --type "SecureString" \
    --region "$AWS_REGION" \
    --profile "$AWS_PROFILE"
fi

echo
echo "✅ All parameters uploaded as version 1 (fresh)."