package com.okblog.post.repository;

import com.okblog.post.model.Post;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface PostRepository extends JpaRepository<Post, UUID> {
    
    List<Post> findByProfileId(UUID profileId);
    
    List<Post> findByProfileIdAndIsPublished(UUID profileId, boolean isPublished);
    
    List<Post> findByType(Post.PostType type);
    
    List<Post> findByTypeAndIsPublished(Post.PostType type, boolean isPublished);
    
    Optional<Post> findBySlug(String slug);
    
    List<Post> findByTagsContaining(String tag);
} 