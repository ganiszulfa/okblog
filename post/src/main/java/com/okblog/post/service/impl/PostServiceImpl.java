package com.okblog.post.service.impl;

import com.okblog.post.dto.PageResponse;
import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;
import com.okblog.post.repository.PostRepository;
import com.okblog.post.service.PostService;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
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
    public PageResponse<PostResponse> createPost(PostRequest request, UUID userId) {
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
        PostResponse postResponse = mapToPostResponse(savedPost);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<PostResponse> getPostById(UUID id) {
        Post post = findPostById(id);
        PostResponse postResponse = mapToPostResponse(post);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<PostResponse> getPostBySlug(String slug) {
        Post post = postRepository.findBySlug(slug)
                .orElseThrow(() -> new EntityNotFoundException("Post not found with slug: " + slug));
        PostResponse postResponse = mapToPostResponse(post);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<List<PostResponse>> getAllPosts(int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Post> postPage = postRepository.findAll(pageable);
        
        List<PostResponse> postResponses = postPage.getContent().stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
        
        return PageResponse.<List<PostResponse>>builder()
                .data(postResponses)
                .pagination(buildPaginationMetadata(postPage, page, perPage))
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<List<PostResponse>> getPostsByProfileId(UUID profileId, int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Post> postPage = postRepository.findByProfileId(profileId, pageable);
        
        List<PostResponse> postResponses = postPage.getContent().stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
        
        return PageResponse.<List<PostResponse>>builder()
                .data(postResponses)
                .pagination(buildPaginationMetadata(postPage, page, perPage))
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<List<PostResponse>> getPostsByProfileIdAndPublished(UUID profileId, boolean isPublished, int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Post> postPage = postRepository.findByProfileIdAndIsPublished(profileId, isPublished, pageable);
        
        List<PostResponse> postResponses = postPage.getContent().stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
        
        return PageResponse.<List<PostResponse>>builder()
                .data(postResponses)
                .pagination(buildPaginationMetadata(postPage, page, perPage))
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<List<PostResponse>> getPostsByType(Post.PostType type, int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Post> postPage = postRepository.findByType(type, pageable);
        
        List<PostResponse> postResponses = postPage.getContent().stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
        
        return PageResponse.<List<PostResponse>>builder()
                .data(postResponses)
                .pagination(buildPaginationMetadata(postPage, page, perPage))
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<List<PostResponse>> getPostsByTypeAndPublished(Post.PostType type, boolean isPublished, int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Post> postPage = postRepository.findByTypeAndIsPublished(type, isPublished, pageable);
        
        List<PostResponse> postResponses = postPage.getContent().stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
        
        return PageResponse.<List<PostResponse>>builder()
                .data(postResponses)
                .pagination(buildPaginationMetadata(postPage, page, perPage))
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<List<PostResponse>> getPostsByTag(String tag, int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Post> postPage = postRepository.findByTagsContaining(tag, pageable);
        
        List<PostResponse> postResponses = postPage.getContent().stream()
                .map(this::mapToPostResponse)
                .collect(Collectors.toList());
        
        return PageResponse.<List<PostResponse>>builder()
                .data(postResponses)
                .pagination(buildPaginationMetadata(postPage, page, perPage))
                .build();
    }
    
    @Override
    @Transactional
    public PageResponse<PostResponse> updatePost(UUID id, PostRequest request) {
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
        PostResponse postResponse = mapToPostResponse(updatedPost);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    @Override
    @Transactional
    public PageResponse<PostResponse> publishPost(UUID id) {
        Post post = findPostById(id);
        post.setPublished(true);
        post.setPublishedAt(LocalDateTime.now());
        post.setUpdatedAt(LocalDateTime.now());
        Post updatedPost = postRepository.save(post);
        PostResponse postResponse = mapToPostResponse(updatedPost);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    @Override
    @Transactional
    public PageResponse<PostResponse> unpublishPost(UUID id) {
        Post post = findPostById(id);
        post.setPublished(false);
        post.setPublishedAt(null);
        post.setUpdatedAt(LocalDateTime.now());
        Post updatedPost = postRepository.save(post);
        PostResponse postResponse = mapToPostResponse(updatedPost);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    @Override
    @Transactional
    public void deletePost(UUID id) {
        Post post = findPostById(id);
        postRepository.delete(post);
    }
    
    @Override
    @Transactional
    public PageResponse<PostResponse> incrementViewCount(UUID id) {
        Post post = findPostById(id);
        post.setViewCount(post.getViewCount() + 1);
        Post updatedPost = postRepository.save(post);
        PostResponse postResponse = mapToPostResponse(updatedPost);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
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
                .publishedAt(post.getPublishedAt())
                .tags(post.getTags())
                .isPublished(post.isPublished())
                .slug(post.getSlug())
                .excerpt(post.getExcerpt())
                .viewCount(post.getViewCount())
                .build();
    }
    
    private PageResponse.PaginationMetadata buildPaginationMetadata(Page<?> page, int currentPage, int perPage) {
        return PageResponse.PaginationMetadata.builder()
                .current_page(currentPage)
                .per_page(perPage)
                .total_pages(page.getTotalPages())
                .total_items(page.getTotalElements())
                .next_page(page.hasNext() ? currentPage + 1 : null)
                .prev_page(page.hasPrevious() ? currentPage - 1 : null)
                .build();
    }
    
    private PageResponse.PaginationMetadata buildSingleItemPagination() {
        return PageResponse.PaginationMetadata.builder()
                .current_page(1)
                .per_page(1)
                .total_pages(1)
                .total_items(1)
                .next_page(null)
                .prev_page(null)
                .build();
    }
} 