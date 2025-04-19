package com.okblog.post.resolver;

import com.okblog.post.annotation.RequiresUserId;
import com.okblog.post.exception.UnauthorizedException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.MethodParameter;
import org.springframework.stereotype.Component;
import org.springframework.web.bind.support.WebDataBinderFactory;
import org.springframework.web.context.request.NativeWebRequest;
import org.springframework.web.method.support.HandlerMethodArgumentResolver;
import org.springframework.web.method.support.ModelAndViewContainer;

import java.util.Base64;
import java.util.UUID;
import java.util.Map;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * Resolver that handles the @RequiresUserId annotation by extracting the user ID from the JWT token
 * in the Authorization header.
 */
@Slf4j
@Component
public class UserIdArgumentResolver implements HandlerMethodArgumentResolver {

    private static final String AUTHORIZATION_HEADER = "Authorization";
    private static final String BEARER_PREFIX = "Bearer ";
    private static final ObjectMapper objectMapper = new ObjectMapper();

    @Override
    public boolean supportsParameter(MethodParameter parameter) {
        return parameter.getParameterAnnotation(RequiresUserId.class) != null
                && parameter.getParameterType().equals(UUID.class);
    }

    @Override
    public Object resolveArgument(MethodParameter parameter, ModelAndViewContainer mavContainer,
                                  NativeWebRequest webRequest, WebDataBinderFactory binderFactory) {
        // Log all headers
        log.info("Received headers:");
        webRequest.getHeaderNames().forEachRemaining(headerName -> {
            String headerValue = webRequest.getHeader(headerName);
            log.info("Header: {} = {}", headerName, headerValue);
        });
        
        String authHeader = webRequest.getHeader(AUTHORIZATION_HEADER);
        log.info("Authorization header value: {}", authHeader);
        
        if (authHeader == null || !authHeader.startsWith(BEARER_PREFIX)) {
            throw new UnauthorizedException("Authorization header with Bearer token is required");
        }
        
        // Extract the JWT token from the Authorization header
        String token = authHeader.substring(BEARER_PREFIX.length());
        
        try {
            // Extract the userId from the JWT claims
            String userId = extractUserIdFromToken(token);
            log.info("Extracted userId from JWT token: {}", userId);
            
            if (userId == null || userId.isBlank()) {
                throw new UnauthorizedException("User ID not found in JWT token");
            }
            
            return UUID.fromString(userId);
        } catch (IllegalArgumentException e) {
            throw new UnauthorizedException("Invalid user ID format in JWT token");
        } catch (Exception e) {
            log.error("Error extracting userId from JWT token", e);
            throw new UnauthorizedException("Error processing JWT token: " + e.getMessage());
        }
    }
    
    /**
     * Extracts the userId claim from a JWT token.
     * 
     * @param token the JWT token
     * @return the userId from the token's payload
     * @throws Exception if token parsing fails
     */
    private String extractUserIdFromToken(String token) throws Exception {
        String[] parts = token.split("\\.");
        if (parts.length != 3) {
            throw new IllegalArgumentException("Invalid JWT token format");
        }
        
        // Decode the payload (second part of the token)
        String payload = new String(Base64.getUrlDecoder().decode(parts[1]));
        Map<String, Object> claims = objectMapper.readValue(payload, Map.class);
        
        // Extract the userId claim
        return (String) claims.get("userId");
    }
} 