FROM eclipse-temurin:21-jdk-alpine as build
WORKDIR /workspace/app

# Install Maven and use --no-cache to avoid additional layer
RUN apk add --no-cache maven

COPY pom.xml .
# Cache maven dependencies
RUN mvn dependency:go-offline || true

COPY src src

# Build the application and clean up Maven cache
RUN mvn package -DskipTests \
    && mkdir -p target/dependency \
    && (cd target/dependency; jar -xf ../*.jar) \
    && rm -rf /root/.m2/repository

# Smaller runtime image
FROM eclipse-temurin:21-jre-alpine
VOLUME /tmp
ARG DEPENDENCY=/workspace/app/target/dependency

# Copy only the necessary files
COPY --from=build ${DEPENDENCY}/BOOT-INF/lib /app/lib
COPY --from=build ${DEPENDENCY}/META-INF /app/META-INF
COPY --from=build ${DEPENDENCY}/BOOT-INF/classes /app

# Use a non-root user for better security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["java","-cp","app:app/lib/*","com.okblog.post.PostServiceApplication"]

EXPOSE 8081 