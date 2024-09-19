docker build -t columbus-crm-backend .
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 767397881306.dkr.ecr.ap-southeast-1.amazonaws.com
docker tag columbus-crm-backend:latest 767397881306.dkr.ecr.ap-southeast-1.amazonaws.com/columbus-crm-backend:latest
docker push 767397881306.dkr.ecr.ap-southeast-1.amazonaws.com/columbus-crm-backend:latest
