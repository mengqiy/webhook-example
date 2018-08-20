#!/usr/bin/env bash

kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io mutating-webhook-configuration

kubectl delete validatingwebhookconfigurations.admissionregistration.k8s.io validating-webhook-configuration

kubectl delete service foo-admission-server-service

kubectl delete deployment wh

kubectl delete deployment nginx
