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

import java.util.UUID;

/**
 * Resolver that handles the @RequiresUserId annotation by extracting the user ID from the X-USERID header.
 */
@Slf4j
@Component
public class UserIdArgumentResolver implements HandlerMethodArgumentResolver {

    private static final String X_USERID_HEADER = "X-USERID";

    @Override
    public boolean supportsParameter(MethodParameter parameter) {
        return parameter.getParameterAnnotation(RequiresUserId.class) != null
                && parameter.getParameterType().equals(UUID.class);
    }

    @Override
    public Object resolveArgument(MethodParameter parameter, ModelAndViewContainer mavContainer,
                                  NativeWebRequest webRequest, WebDataBinderFactory binderFactory) {
        String userId = webRequest.getHeader(X_USERID_HEADER);

        log.info("UserIdArgumentResolver - X-USERID header value: {}", userId);
        
        if (userId == null || userId.isBlank()) {
            throw new UnauthorizedException("User ID is required");
        }
        
        try {
            return UUID.fromString(userId);
        } catch (IllegalArgumentException e) {
            throw new UnauthorizedException("Invalid user ID format");
        }
    }
} 