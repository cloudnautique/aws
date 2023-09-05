FROM public.ecr.aws/lambda/nodejs:18 as dev
COPY . .
RUN npm install && \
    npx tsc 
CMD [ "built/app.handler" ]

FROM public.ecr.aws/lambda/nodejs:18 as prod
COPY --from=dev /var/task/built/*.js ./
COPY --from=dev /var/task/package*.json ./
RUN npm install
CMD [ "app.handler" ]