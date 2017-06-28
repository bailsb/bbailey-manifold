# bbailey-manifold
manifold project for ben

uses goji network framework (https://github.com/zenazn/goji)

to test;

# create a test master key for grafton to use when acting as Manifold
# this file is written as masterkey.json
grafton generate

# Set environment variables to configure the test app
# MASTER_KEY is the public_key portion of masterkey.json
export MASTER_KEY="PUBLIC-KEY-FROM-GRAFTON"

# CONNECTOR_URL is the url that Grafton will listen on. It corresponds to
# Grafton's --sso-port flag.
export CONNECTOR_URL=http://localhost:3001/v1

# Set fake OAuth 2.0 credentials. The format of these are specific, so you can
# reuse the values here.
export CLIENT_ID=21jtaatqj8y5t0kctb2ejr6jev5w8
export CLIENT_SECRET=3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw

then go run app.go

In another terminal run grafton test with the following;

rafton test --product=numbers --plan=small --region=aws::us-east-1 \
    --client-id=21jtaatqj8y5t0kctb2ejr6jev5w8 \
    --client-secret=3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw \
    --connector-port=3001 \
    --new-plan=large \
    http://localhost:8000
