<h1 align="center">Welcome to razorpay-go-lambda 👋</h1>
<p>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> Lambda function for Razorpay Subscriptions and E-mandates.

Note: Switch over to the `erpnext_integration` branch for ERPNext integration.

## Install

Serverless framework setup for AWS is a pre-requirement.

```sh
make build
```

## Usage

Set the following env variables before depoying:

- `KEY` = razorpay key
- `SECRET` = razorpay secret

If you're on the `erpnext_integration` branch set these additional env variables:

- `ERP_KEY` = ERPNext key
- `ERP_SECRET` = ERPNext secret

```sh
make deploy
```

## Author

👤 **Anoop B**

- Github: [@anoop-b](https://github.com/anoop-b)

## Show your support

Give a ⭐️ if this project helped you!
