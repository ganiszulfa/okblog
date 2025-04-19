package com.okblog.post.dto;

import com.okblog.post.model.Post;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Set;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class PostRequest {
    
    @NotNull
    private Post.PostType type;
    
    @NotBlank
    private String title;
    
    @NotBlank
    private String content;
    
    private Set<String> tags;
    
    private boolean isPublished;
    
    private String slug;
    
    private String excerpt;
} 