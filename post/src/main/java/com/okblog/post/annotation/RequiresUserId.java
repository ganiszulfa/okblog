package com.okblog.post.annotation;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Annotation that indicates a parameter should be automatically filled with 
 * the user ID from the X-USERID header.
 * If the header is missing or invalid, a 401 Unauthorized response will be returned.
 * 
 * This annotation can be used on method parameters of type UUID.
 */
@Target(ElementType.PARAMETER)
@Retention(RetentionPolicy.RUNTIME)
public @interface RequiresUserId {
} 