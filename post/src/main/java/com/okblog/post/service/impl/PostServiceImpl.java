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
import org.springframework.data.domain.Sort;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
public class PostServiceImpl implements PostService {

    private final PostRepository postRepository;
    
    // In-memory cache for view counts, mapping post ID to view counter
    private final Map<UUID, AtomicInteger> viewCountCache = new ConcurrentHashMap<>();
    
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
    public PageResponse<PostResponse> getPostBySlug(String slug, boolean onlyPublished) {
        Post post = postRepository.findBySlug(slug)
                .orElseThrow(() -> new EntityNotFoundException("Post not found with slug: " + slug));
                
        if (onlyPublished && !post.isPublished()) {
            throw new EntityNotFoundException("Published post not found with slug: " + slug);
        }
        
        PostResponse postResponse = mapToPostResponse(post);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    @Override
    @Transactional(readOnly = true)
    public PageResponse<List<PostResponse>> getAllPosts(int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage, Sort.by(Sort.Direction.DESC, "publishedAt"));
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
        Pageable pageable = PageRequest.of(page - 1, perPage, Sort.by(Sort.Direction.DESC, "publishedAt"));
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
        Pageable pageable = PageRequest.of(page - 1, perPage, Sort.by(Sort.Direction.DESC, "publishedAt"));
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
        Pageable pageable = PageRequest.of(page - 1, perPage, Sort.by(Sort.Direction.DESC, "publishedAt"));
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
        Pageable pageable = PageRequest.of(page - 1, perPage, Sort.by(Sort.Direction.DESC, "publishedAt"));
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
        Pageable pageable = PageRequest.of(page - 1, perPage, Sort.by(Sort.Direction.DESC, "publishedAt"));
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
        
        // Get or create the atomic counter for this post
        AtomicInteger counter = viewCountCache.computeIfAbsent(id, k -> new AtomicInteger(0));
        
        // Increment the counter
        int currentCount = counter.incrementAndGet();
        
        // Return post with the updated view count (from cache + database)
        PostResponse postResponse = mapToPostResponse(post);
        // Add the count from cache that hasn't been persisted yet
        postResponse.setViewCount(postResponse.getViewCount() + currentCount);
        
        return PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(buildSingleItemPagination())
                .build();
    }
    
    /**
     * Flushes the cached view count for a post to the database
     */
    private void flushViewCount(UUID postId, Post post) {
        AtomicInteger counter = viewCountCache.get(postId);
        if (counter != null) {
            int count = counter.getAndSet(0); // Reset counter
            if (count > 0) {
                post.setViewCount(post.getViewCount() + count);
                postRepository.save(post);
                log.debug("Flushed {} views for post {}", count, postId);
            }
        }
    }
    
    /**
     * Scheduled job to flush all view counts to the database
     * Runs every 5 minutes
     */
    @Scheduled(fixedRate = 5 * 60 * 1000) 
    @Transactional
    public void flushAllViewCounts() {
        log.info("Starting scheduled view count flush");
        
        // Make a copy of the keys to avoid concurrent modification
        List<UUID> postIds = List.copyOf(viewCountCache.keySet());
        
        for (UUID postId : postIds) {
            try {
                Post post = postRepository.findById(postId).orElse(null);
                if (post != null) {
                    flushViewCount(postId, post);
                } else {
                    // Post was deleted, remove from cache
                    viewCountCache.remove(postId);
                }
            } catch (Exception e) {
                log.error("Error flushing view count for post {}", postId, e);
            }
        }
        
        log.info("Completed scheduled view count flush");
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