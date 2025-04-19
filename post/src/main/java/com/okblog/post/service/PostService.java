package com.okblog.post.service;

import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;

import java.util.List;
import java.util.UUID;

public interface PostService {
    
    PostResponse createPost(PostRequest request, UUID userId);
    
    PostResponse getPostById(UUID id);
    
    PostResponse getPostBySlug(String slug);
    
    List<PostResponse> getAllPosts();
    
    List<PostResponse> getPostsByProfileId(UUID profileId);
    
    List<PostResponse> getPostsByProfileIdAndPublished(UUID profileId, boolean isPublished);
    
    List<PostResponse> getPostsByType(Post.PostType type);
    
    List<PostResponse> getPostsByTypeAndPublished(Post.PostType type, boolean isPublished);
    
    List<PostResponse> getPostsByTag(String tag);
    
    PostResponse updatePost(UUID id, PostRequest request);
    
    PostResponse publishPost(UUID id);
    
    PostResponse unpublishPost(UUID id);
    
    void deletePost(UUID id);
    
    PostResponse incrementViewCount(UUID id);
} 