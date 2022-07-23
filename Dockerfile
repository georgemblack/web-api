FROM node:18 AS build-env
ADD . /build
WORKDIR /build
RUN yarn && yarn build

FROM node:18-alpine
ENV NODE_ENV=production
COPY --from=build-env /build /app
WORKDIR /app
CMD ["yarn", "start"]
