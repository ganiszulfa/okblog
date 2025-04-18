package com.okblog.post.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;
import com.okblog.post.service.PostService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Set;
import java.util.UUID;

import static org.hamcrest.Matchers.hasSize;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(PostController.class)
public class PostControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private PostService postService;

    @Autowired
    private ObjectMapper objectMapper;

    private UUID postId;
    private UUID profileId;
    private PostResponse postResponse;
    private PostRequest postRequest;

    @BeforeEach
    void setUp() {
        postId = UUID.randomUUID();
        profileId = UUID.randomUUID();

        postRequest = PostRequest.builder()
                .profileId(profileId)
                .type(Post.PostType.POST)
                .title("Test Post")
                .content("Test content")
                .tags(Set.of("test", "unit-test"))
                .isPublished(true)
                .slug("test-post")
                .excerpt("Test excerpt")
                .build();

        postResponse = PostResponse.builder()
                .id(postId)
                .profileId(profileId)
                .type(Post.PostType.POST)
                .title("Test Post")
                .content("Test content")
                .createdAt(LocalDateTime.now())
                .tags(Set.of("test", "unit-test"))
                .isPublished(true)
                .slug("test-post")
                .excerpt("Test excerpt")
                .viewCount(0)
                .build();
    }

    @Test
    void whenCreatePost_thenReturn201AndPost() throws Exception {
        when(postService.createPost(any(PostRequest.class))).thenReturn(postResponse);

        mockMvc.perform(post("/api/posts")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(postRequest)))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").value(postId.toString()))
                .andExpect(jsonPath("$.title").value("Test Post"));
    }

    @Test
    void whenGetAllPosts_thenReturn200AndPostList() throws Exception {
        when(postService.getAllPosts()).thenReturn(List.of(postResponse));

        mockMvc.perform(get("/api/posts"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(1)))
                .andExpect(jsonPath("$[0].id").value(postId.toString()));
    }

    @Test
    void whenGetPostById_thenReturn200AndPost() throws Exception {
        when(postService.getPostById(postId)).thenReturn(postResponse);

        mockMvc.perform(get("/api/posts/{id}", postId))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(postId.toString()))
                .andExpect(jsonPath("$.title").value("Test Post"));
    }
} 