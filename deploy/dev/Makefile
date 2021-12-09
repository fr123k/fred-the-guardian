export KUBECONFIG=$(shell pwd)/kubeconfig.yaml

deploy:
	kubectl apply -f deployment.yml 
	kubectl apply -f service.yml
	kubectl apply -f ingress.yml
