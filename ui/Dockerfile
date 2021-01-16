FROM cypress/included:4.12.1
# replace this with your application's default port
WORKDIR /work
COPY yarn.lock .
COPY package.json .
RUN yarn

COPY . .
ENV CYPRESS_BASE_URL=http://localhost:3001
ENV REACT_APP_API_URL=http://gourd:4242
EXPOSE 3001
RUN pwd
RUN ls
ENTRYPOINT [ "./ci.sh" ]