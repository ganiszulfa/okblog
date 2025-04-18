package com.okblog.post.controller;

import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;
import com.okblog.post.service.PostService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/posts")
@RequiredArgsConstructor
public class PostController {

    private final PostService postService;
    
    @PostMapping
    public ResponseEntity<PostResponse> createPost(@Valid @RequestBody PostRequest request) {
        return new ResponseEntity<>(postService.createPost(request), HttpStatus.CREATED);
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<PostResponse> getPostById(@PathVariable UUID id) {
        return ResponseEntity.ok(postService.getPostById(id));
    }
    
    @GetMapping("/slug/{slug}")
    public ResponseEntity<PostResponse> getPostBySlug(@PathVariable String slug) {
        return ResponseEntity.ok(postService.getPostBySlug(slug));
    }
    
    @GetMapping
    public ResponseEntity<List<PostResponse>> getAllPosts() {
        return ResponseEntity.ok(postService.getAllPosts());
    }
    
    @GetMapping("/profile/{profileId}")
    public ResponseEntity<List<PostResponse>> getPostsByProfileId(@PathVariable UUID profileId) {
        return ResponseEntity.ok(postService.getPostsByProfileId(profileId));
    }
    
    @GetMapping("/profile/{profileId}/published/{isPublished}")
    public ResponseEntity<List<PostResponse>> getPostsByProfileIdAndPublished(
            @PathVariable UUID profileId, 
            @PathVariable boolean isPublished) {
        return ResponseEntity.ok(postService.getPostsByProfileIdAndPublished(profileId, isPublished));
    }
    
    @GetMapping("/type/{type}")
    public ResponseEntity<List<PostResponse>> getPostsByType(@PathVariable Post.PostType type) {
        return ResponseEntity.ok(postService.getPostsByType(type));
    }
    
    @GetMapping("/type/{type}/published/{isPublished}")
    public ResponseEntity<List<PostResponse>> getPostsByTypeAndPublished(
            @PathVariable Post.PostType type, 
            @PathVariable boolean isPublished) {
        return ResponseEntity.ok(postService.getPostsByTypeAndPublished(type, isPublished));
    }
    
    @GetMapping("/tag/{tag}")
    public ResponseEntity<List<PostResponse>> getPostsByTag(@PathVariable String tag) {
        return ResponseEntity.ok(postService.getPostsByTag(tag));
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<PostResponse> updatePost(
            @PathVariable UUID id, 
            @Valid @RequestBody PostRequest request) {
        return ResponseEntity.ok(postService.updatePost(id, request));
    }
    
    @PutMapping("/{id}/publish")
    public ResponseEntity<PostResponse> publishPost(@PathVariable UUID id) {
        return ResponseEntity.ok(postService.publishPost(id));
    }
    
    @PutMapping("/{id}/unpublish")
    public ResponseEntity<PostResponse> unpublishPost(@PathVariable UUID id) {
        return ResponseEntity.ok(postService.unpublishPost(id));
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deletePost(@PathVariable UUID id) {
        postService.deletePost(id);
        return ResponseEntity.noContent().build();
    }
    
    @PutMapping("/{id}/view")
    public ResponseEntity<PostResponse> incrementViewCount(@PathVariable UUID id) {
        return ResponseEntity.ok(postService.incrementViewCount(id));
    }
} 