package com.okblog.post.repository;

import com.okblog.post.model.Post;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface PostRepository extends JpaRepository<Post, UUID> {
    
    Page<Post> findByProfileId(UUID profileId, Pageable pageable);
    
    Page<Post> findByProfileIdAndIsPublished(UUID profileId, boolean isPublished, Pageable pageable);
    
    Page<Post> findByType(Post.PostType type, Pageable pageable);
    
    Page<Post> findByTypeAndIsPublished(Post.PostType type, boolean isPublished, Pageable pageable);
    
    Optional<Post> findBySlug(String slug);
    
    Page<Post> findByTagsContaining(String tag, Pageable pageable);
} 