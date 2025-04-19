package com.okblog.post.logging;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.stereotype.Component;

import java.time.Duration;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

@Slf4j
@Aspect
@Component
@RequiredArgsConstructor
public class KibanaLoggerAspect {

    private final KibanaLogger kibanaLogger;

    /**
     * Log all service method executions to Kibana
     */
    @Around("execution(* com.okblog.post.service.*.*(..))")
    public Object logServiceMethodExecutionTime(ProceedingJoinPoint joinPoint) throws Throwable {
        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        String methodName = signature.getMethod().getName();
        String className = signature.getDeclaringType().getSimpleName();

        Instant start = Instant.now();
        Object result = null;
        Throwable thrownException = null;

        try {
            // Execute the actual method
            result = joinPoint.proceed();
            return result;
        } catch (Throwable exception) {
            thrownException = exception;
            throw exception;
        } finally {
            // Calculate execution time
            Duration duration = Duration.between(start, Instant.now());
            
            // Build the log fields
            Map<String, Object> fields = new HashMap<>();
            fields.put("class", className);
            fields.put("method", methodName);
            fields.put("duration_ms", duration.toMillis());
            fields.put("execution_id", UUID.randomUUID().toString());
            
            // Add parameter names and values
            String[] paramNames = signature.getParameterNames();
            Object[] paramValues = joinPoint.getArgs();
            Map<String, Object> params = new HashMap<>();
            
            for (int i = 0; i < paramNames.length; i++) {
                // For security/privacy reasons, don't log the actual content of sensitive objects
                if (paramValues[i] instanceof UUID) {
                    params.put(paramNames[i], paramValues[i].toString());
                } else if (paramValues[i] != null) {
                    params.put(paramNames[i], paramValues[i].getClass().getSimpleName() + "@" + 
                               System.identityHashCode(paramValues[i]));
                } else {
                    params.put(paramNames[i], "null");
                }
            }
            fields.put("parameters", params);
            
            // Add result info
            if (thrownException != null) {
                fields.put("status", "ERROR");
                fields.put("exception", thrownException.getClass().getName());
                fields.put("error_message", thrownException.getMessage());
                
                kibanaLogger.error(
                    String.format("Method %s.%s failed in %d ms", className, methodName, duration.toMillis()),
                    fields
                );
            } else {
                fields.put("status", "SUCCESS");
                
                if (result != null) {
                    fields.put("result_type", result.getClass().getName());
                } else {
                    fields.put("result_type", "void/null");
                }
                
                kibanaLogger.info(
                    String.format("Method %s.%s executed in %d ms", className, methodName, duration.toMillis()),
                    fields
                );
            }
        }
    }
} 