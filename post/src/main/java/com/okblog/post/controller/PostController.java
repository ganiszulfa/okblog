package com.okblog.post.controller;

import com.okblog.post.annotation.RequiresUserId;
import com.okblog.post.dto.PageResponse;
import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;
import com.okblog.post.service.PostService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/posts")
@RequiredArgsConstructor
@Slf4j
public class PostController {

    private final PostService postService;
    private static final int DEFAULT_PAGE = 1;
    private static final int DEFAULT_PER_PAGE = 10;
    
    @PostMapping
    public ResponseEntity<PageResponse<PostResponse>> createPost(@Valid @RequestBody PostRequest request, @RequiresUserId UUID userId) {
        log.info("Authenticated userId from JWT token: {}", userId);
        return new ResponseEntity<>(postService.createPost(request, userId), HttpStatus.CREATED);
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<PageResponse<PostResponse>> getPostById(@PathVariable UUID id, @RequiresUserId UUID userId) {
        return ResponseEntity.ok(postService.getPostById(id));
    }
    
    @GetMapping("/slug/{slug}")
    public ResponseEntity<PageResponse<PostResponse>> getPostBySlug(@PathVariable String slug) {
        return ResponseEntity.ok(postService.getPostBySlug(slug, true));
    }
    
    @GetMapping
    public ResponseEntity<PageResponse<List<PostResponse>>> getAllPosts(
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int per_page) {
        return ResponseEntity.ok(postService.getAllPosts(page, per_page));
    }
    
    @GetMapping("/profile/{profileId}")
    public ResponseEntity<PageResponse<List<PostResponse>>> getPostsByProfileId(
            @PathVariable UUID profileId,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int per_page) {
        return ResponseEntity.ok(postService.getPostsByProfileId(profileId, page, per_page));
    }
    
    @GetMapping("/my-posts")
    public ResponseEntity<PageResponse<List<PostResponse>>> getMyPosts(
            @RequiresUserId UUID userId,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int per_page) {
        return ResponseEntity.ok(postService.getPostsByProfileId(userId, page, per_page));
    }
    
    @GetMapping("/my-posts/published/{isPublished}")
    public ResponseEntity<PageResponse<List<PostResponse>>> getMyPostsByPublishedStatus(
            @RequiresUserId UUID userId,
            @PathVariable boolean isPublished,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int per_page) {
        return ResponseEntity.ok(postService.getPostsByProfileIdAndPublished(userId, isPublished, page, per_page));
    }
    
    @GetMapping("/type/{type}")
    public ResponseEntity<PageResponse<List<PostResponse>>> getPostsByType(
            @PathVariable Post.PostType type,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int per_page) {
        return ResponseEntity.ok(postService.getPostsByType(type, page, per_page));
    }
    
    @GetMapping("/type/{type}/published/{isPublished}")
    public ResponseEntity<PageResponse<List<PostResponse>>> getPostsByTypeAndPublished(
            @PathVariable Post.PostType type,
            @PathVariable boolean isPublished,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int per_page) {
        return ResponseEntity.ok(postService.getPostsByTypeAndPublished(type, isPublished, page, per_page));
    }
    
    @GetMapping("/tag/{tag}")
    public ResponseEntity<PageResponse<List<PostResponse>>> getPostsByTag(
            @PathVariable String tag,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int per_page) {
        return ResponseEntity.ok(postService.getPostsByTag(tag, page, per_page));
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<PageResponse<PostResponse>> updatePost(
            @PathVariable UUID id,
            @Valid @RequestBody PostRequest request,
            @RequiresUserId UUID userId) {
        // First check if the post belongs to this user
        ResponseEntity<?> ownershipCheck = checkPostOwnership(id, userId);
        if (ownershipCheck != null) {
            return (ResponseEntity<PageResponse<PostResponse>>) ownershipCheck;
        }
        return ResponseEntity.ok(postService.updatePost(id, request));
    }
    
    @PutMapping("/{id}/publish")
    public ResponseEntity<PageResponse<PostResponse>> publishPost(@PathVariable UUID id, @RequiresUserId UUID userId) {
        // First check if the post belongs to this user
        ResponseEntity<?> ownershipCheck = checkPostOwnership(id, userId);
        if (ownershipCheck != null) {
            return (ResponseEntity<PageResponse<PostResponse>>) ownershipCheck;
        }
        return ResponseEntity.ok(postService.publishPost(id));
    }
    
    @PutMapping("/{id}/unpublish")
    public ResponseEntity<PageResponse<PostResponse>> unpublishPost(@PathVariable UUID id, @RequiresUserId UUID userId) {
        // First check if the post belongs to this user
        ResponseEntity<?> ownershipCheck = checkPostOwnership(id, userId);
        if (ownershipCheck != null) {
            return (ResponseEntity<PageResponse<PostResponse>>) ownershipCheck;
        }
        return ResponseEntity.ok(postService.unpublishPost(id));
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deletePost(@PathVariable UUID id, @RequiresUserId UUID userId) {
        // First check if the post belongs to this user
        ResponseEntity<?> ownershipCheck = checkPostOwnership(id, userId);
        if (ownershipCheck != null) {
            return (ResponseEntity<Void>) ownershipCheck;
        }
        postService.deletePost(id);
        return ResponseEntity.noContent().build();
    }
    
    @PutMapping("/{id}/view")
    public ResponseEntity<PageResponse<PostResponse>> incrementViewCount(@PathVariable UUID id) {
        return ResponseEntity.ok(postService.incrementViewCount(id));
    }
    

    private ResponseEntity<?> checkPostOwnership(UUID postId, UUID userId) {
        PageResponse<PostResponse> existingPostResponse = postService.getPostById(postId);
        PostResponse existingPost = existingPostResponse.getData();
        if (!existingPost.getProfileId().equals(userId)) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).build();
        }
        return null;
    }
    
} 