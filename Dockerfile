FROM node:13-stretch AS build-env
ADD . /build
WORKDIR /build
RUN yarn

FROM node:13-alpine
COPY --from=build-env /build /app
WORKDIR /app
CMD ["index.js"]