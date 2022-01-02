export PING_ADDRESS=$(shell minikube service fred-the-guardian --url)
export EXTERNAL_IP?=172.28.128.16

vagrant:
	vagrant up --provider virtualbox 

setup:
	# There is an issue with the nginx validation webhook in minikube make it impossible to deploy ingress resources.
	# This is a temporarly fix
	kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission || echo "ignore errors"

deploy-redis:
	kubectl apply -f redis/deployment.yml 
	kubectl apply -f redis/service.yml

deploy: deploy-redis
	kubectl apply -f deployment.yml 
	kubectl apply -f service.yml
	kubectl apply -f ingress.yml

test-minikube:
	curl -X POST -H 'X-SECRET-KEY:top secret' -v ${PING_ADDRESS}/ping -d '{"request":"ping"}'

test:
	curl -b .cookies -c .cookies -X POST -H 'X-SECRET-KEY:top secret' -v http://${EXTERNAL_IP}/fred/ping -d '{"request":"ping"}'

clean:
	vagrant destroy -f
	rm -f ca.crt client.crt client.key
