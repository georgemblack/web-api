FROM node:13-stretch AS build-env
ADD . /build
WORKDIR /build
RUN yarn

FROM node:13-alpine
ENV NODE_ENV=production
COPY --from=build-env /build /app
WORKDIR /app
CMD ["yarn", "start"]
