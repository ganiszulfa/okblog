package com.okblog.post.service.impl;

import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;
import com.okblog.post.repository.PostRepository;
import com.okblog.post.service.PostService;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
public class PostServiceImpl implements PostService {

    private final PostRepository postRepository;
    
    @Override
    @Transactional
    public PostResponse createPost(PostRequest request, UUID userId) {
        log.info("Authenticated userId from JWT token: {}", userId);
        Post post = Post.builder()
                .profileId(userId)
                .type(request.getType())
                .title(request.getTitle())
                .content(request.getContent())
                .tags(request.getTags())
                .isPublished(request.isPublished())
                .slug(request.getSlug())
                .excerpt(request.getExcerpt())
                .build();
        
        Post savedPost = postRepository.save(post);
        return mapToPostResponse(savedPost);
    }
    
    @Override
    @Transactional(readOnly = true)
    public PostResponse getPostById(UUID id) {
        Post post = findPostById(id);
        return mapToPostResponse(post);
    }
    
    @Override
    @Transactional(readOnly = true)
    public PostResponse getPostBySlug(String slug) {
        Post post = postRepository.findBySlug(slug)
                .orElseThrow(() -> new EntityNotFoundException("Post not found with slug: " + slug));
        return mapToPostResponse(post);
    }
    
    @Override
    @Transactional(readOnly = true)
    public List<PostResponse> getAllPosts() {
        return postRepository.findAll().stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
    }
    
    @Override
    @Transactional(readOnly = true)
    public List<PostResponse> getPostsByProfileId(UUID profileId) {
        return postRepository.findByProfileId(profileId).stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
    }
    
    @Override
    @Transactional(readOnly = true)
    public List<PostResponse> getPostsByProfileIdAndPublished(UUID profileId, boolean isPublished) {
        return postRepository.findByProfileIdAndIsPublished(profileId, isPublished).stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
    }
    
    @Override
    @Transactional(readOnly = true)
    public List<PostResponse> getPostsByType(Post.PostType type) {
        return postRepository.findByType(type).stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
    }
    
    @Override
    @Transactional(readOnly = true)
    public List<PostResponse> getPostsByTypeAndPublished(Post.PostType type, boolean isPublished) {
        return postRepository.findByTypeAndIsPublished(type, isPublished).stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
    }
    
    @Override
    @Transactional(readOnly = true)
    public List<PostResponse> getPostsByTag(String tag) {
        return postRepository.findByTagsContaining(tag).stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
    }
    
    @Override
    @Transactional
    public PostResponse updatePost(UUID id, PostRequest request) {
        Post post = findPostById(id);
        
        // post.setProfileId(request.getProfileId());
        post.setType(request.getType());
        post.setTitle(request.getTitle());
        post.setContent(request.getContent());
        post.setTags(request.getTags());
        post.setPublished(request.isPublished());
        post.setSlug(request.getSlug());
        post.setExcerpt(request.getExcerpt());
        post.setUpdatedAt(LocalDateTime.now());
        
        Post updatedPost = postRepository.save(post);
        return mapToPostResponse(updatedPost);
    }
    
    @Override
    @Transactional
    public PostResponse publishPost(UUID id) {
        Post post = findPostById(id);
        post.setPublished(true);
        post.setUpdatedAt(LocalDateTime.now());
        Post updatedPost = postRepository.save(post);
        return mapToPostResponse(updatedPost);
    }
    
    @Override
    @Transactional
    public PostResponse unpublishPost(UUID id) {
        Post post = findPostById(id);
        post.setPublished(false);
        post.setUpdatedAt(LocalDateTime.now());
        Post updatedPost = postRepository.save(post);
        return mapToPostResponse(updatedPost);
    }
    
    @Override
    @Transactional
    public void deletePost(UUID id) {
        Post post = findPostById(id);
        postRepository.delete(post);
    }
    
    @Override
    @Transactional
    public PostResponse incrementViewCount(UUID id) {
        Post post = findPostById(id);
        post.setViewCount(post.getViewCount() + 1);
        Post updatedPost = postRepository.save(post);
        return mapToPostResponse(updatedPost);
    }
    
    private Post findPostById(UUID id) {
        return postRepository.findById(id)
                .orElseThrow(() -> new EntityNotFoundException("Post not found with id: " + id));
    }
    
    private PostResponse mapToPostResponse(Post post) {
        return PostResponse.builder()
                .id(post.getId())
                .profileId(post.getProfileId())
                .type(post.getType())
                .title(post.getTitle())
                .content(post.getContent())
                .createdAt(post.getCreatedAt())
                .updatedAt(post.getUpdatedAt())
                .tags(post.getTags())
                .isPublished(post.isPublished())
                .slug(post.getSlug())
                .excerpt(post.getExcerpt())
                .viewCount(post.getViewCount())
                .build();
    }
} 