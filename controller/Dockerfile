# Use java 18 maven
FROM maven:3.8.5-openjdk-18-slim AS build
ADD ./WebController /app
WORKDIR /app
# install dependencies and package
RUN mvn package


# Build (java 17)
FROM openjdk:17-alpine
COPY --from=build /app/target/WebController-1.0-SNAPSHOT.jar /app/WebController-1.0-SNAPSHOT.jar
# Server runs on 8080
EXPOSE 8080
ENTRYPOINT ["java", "-jar","/app/WebController-1.0-SNAPSHOT.jar"]
