package com.okblog.post.service;

import com.okblog.post.dto.PageResponse;
import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;

import java.util.List;
import java.util.UUID;

public interface PostService {
    
    PageResponse<PostResponse> createPost(PostRequest request, UUID userId);
    
    PageResponse<PostResponse> getPostById(UUID id);
    
    PageResponse<PostResponse> getPostBySlug(String slug);
    
    PageResponse<List<PostResponse>> getAllPosts(int page, int perPage);
    
    PageResponse<List<PostResponse>> getPostsByProfileId(UUID profileId, int page, int perPage);
    
    PageResponse<List<PostResponse>> getPostsByProfileIdAndPublished(UUID profileId, boolean isPublished, int page, int perPage);
    
    PageResponse<List<PostResponse>> getPostsByType(Post.PostType type, int page, int perPage);
    
    PageResponse<List<PostResponse>> getPostsByTypeAndPublished(Post.PostType type, boolean isPublished, int page, int perPage);
    
    PageResponse<List<PostResponse>> getPostsByTag(String tag, int page, int perPage);
    
    PageResponse<PostResponse> updatePost(UUID id, PostRequest request);
    
    PageResponse<PostResponse> publishPost(UUID id);
    
    PageResponse<PostResponse> unpublishPost(UUID id);
    
    void deletePost(UUID id);
    
    PageResponse<PostResponse> incrementViewCount(UUID id);
} 