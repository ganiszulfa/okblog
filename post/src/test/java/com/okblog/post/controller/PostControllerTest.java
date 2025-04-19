package com.okblog.post.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.okblog.post.annotation.RequiresUserId;
import com.okblog.post.dto.PostRequest;
import com.okblog.post.dto.PostResponse;
import com.okblog.post.model.Post;
import com.okblog.post.resolver.UserIdArgumentResolver;
import com.okblog.post.service.PostService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.context.annotation.Import;
import org.springframework.core.MethodParameter;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.web.bind.support.WebDataBinderFactory;
import org.springframework.web.context.request.NativeWebRequest;
import org.springframework.web.method.support.ModelAndViewContainer;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Set;
import java.util.UUID;

import static org.hamcrest.Matchers.hasSize;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(PostController.class)
@Import(PostControllerTest.MockUserIdArgumentResolver.class)
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

    // Create a mock resolver that always returns a test UUID for @RequiresUserId parameters
    public static class MockUserIdArgumentResolver extends UserIdArgumentResolver {
        private static final UUID TEST_USER_ID = UUID.fromString("00000000-0000-0000-0000-000000000001");

        @Override
        public boolean supportsParameter(MethodParameter parameter) {
            return parameter.getParameterAnnotation(RequiresUserId.class) != null;
        }

        @Override
        public Object resolveArgument(MethodParameter parameter, ModelAndViewContainer mavContainer,
                                     NativeWebRequest webRequest, WebDataBinderFactory binderFactory) {
            return TEST_USER_ID;
        }
    }

    @BeforeEach
    void setUp() {
        postId = UUID.randomUUID();
        profileId = MockUserIdArgumentResolver.TEST_USER_ID; // Use same UUID for tests

        postRequest = PostRequest.builder()
                // .profileId(profileId)
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
        // Use any UUID for the second parameter
        when(postService.createPost(any(PostRequest.class), any(UUID.class))).thenReturn(postResponse);

        mockMvc.perform(post("/api/posts")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(postRequest))
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString())) // Add the header
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").value(postId.toString()))
                .andExpect(jsonPath("$.profileId").value(profileId.toString()))
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
    
    @Test
    void whenGetMyPosts_thenReturn200AndPostList() throws Exception {
        when(postService.getPostsByProfileId(MockUserIdArgumentResolver.TEST_USER_ID))
                .thenReturn(List.of(postResponse));

        mockMvc.perform(get("/api/posts/my-posts")
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString()))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(1)))
                .andExpect(jsonPath("$[0].id").value(postId.toString()));
    }
    
    @Test
    void whenUpdatePost_thenReturn200AndUpdatedPost() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(postResponse);
        // Mock the updatePost call
        when(postService.updatePost(eq(postId), any(PostRequest.class))).thenReturn(postResponse);

        mockMvc.perform(put("/api/posts/{id}", postId)
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(postRequest))
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString()))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(postId.toString()));
    }
    
    @Test
    void whenUpdatePostWithUnauthorizedUser_thenReturn403() throws Exception {
        // Create a response with different profileId
        PostResponse unauthorizedResponse = PostResponse.builder()
                .id(postId)
                .profileId(UUID.randomUUID()) // Different from TEST_USER_ID
                .type(Post.PostType.POST)
                .title("Test Post")
                .content("Test content")
                .build();
                
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(unauthorizedResponse);

        mockMvc.perform(put("/api/posts/{id}", postId)
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(postRequest))
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString()))
                .andExpect(status().isForbidden());
    }
    
    @Test
    void whenPublishPost_thenReturn200AndPublishedPost() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(postResponse);
        // Mock the publishPost call
        when(postService.publishPost(postId)).thenReturn(postResponse);

        mockMvc.perform(put("/api/posts/{id}/publish", postId)
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString()))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(postId.toString()));
    }
    
    @Test
    void whenUnpublishPost_thenReturn200AndUnpublishedPost() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(postResponse);
        // Mock the unpublishPost call
        when(postService.unpublishPost(postId)).thenReturn(postResponse);

        mockMvc.perform(put("/api/posts/{id}/unpublish", postId)
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString()))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(postId.toString()));
    }
    
    @Test
    void whenDeletePost_thenReturn204() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(postResponse);
        // Mock the deletePost call
        doNothing().when(postService).deletePost(postId);

        mockMvc.perform(delete("/api/posts/{id}", postId)
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString()))
                .andExpect(status().isNoContent());
    }
    
    @Test
    void whenDeletePostWithUnauthorizedUser_thenReturn403() throws Exception {
        // Create a response with different profileId
        PostResponse unauthorizedResponse = PostResponse.builder()
                .id(postId)
                .profileId(UUID.randomUUID()) // Different from TEST_USER_ID
                .type(Post.PostType.POST)
                .title("Test Post")
                .content("Test content")
                .build();
                
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(unauthorizedResponse);

        mockMvc.perform(delete("/api/posts/{id}", postId)
                .header("X-USERID", MockUserIdArgumentResolver.TEST_USER_ID.toString()))
                .andExpect(status().isForbidden());
    }
} 