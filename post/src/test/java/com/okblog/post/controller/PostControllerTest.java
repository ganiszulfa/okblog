package com.okblog.post.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.okblog.post.annotation.RequiresUserId;
import com.okblog.post.dto.PageResponse;
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
    private PageResponse<PostResponse> singlePostResponse;
    private PageResponse<List<PostResponse>> postListResponse;
    private PostRequest postRequest;
    private static final String TEST_JWT_TOKEN = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDEiLCJ1c2VybmFtZSI6InRlc3RfdXNlciIsImlzc3VlZEF0IjoiMjAyNC0wMS0wMVQwMDowMDowMFoiLCJleHBpcmVzQXQiOiIyMDI0LTAxLTAyVDAwOjAwOjAwWiJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";

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
                // .profileId(profileId) - Removed as profileId is now extracted from JWT token
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
                
        // Create PageResponse for single post
        PageResponse.PaginationMetadata singleItemPagination = PageResponse.PaginationMetadata.builder()
                .current_page(1)
                .per_page(1)
                .total_pages(1)
                .total_items(1)
                .next_page(null)
                .prev_page(null)
                .build();
                
        singlePostResponse = PageResponse.<PostResponse>builder()
                .data(postResponse)
                .pagination(singleItemPagination)
                .build();
                
        // Create PageResponse for post list
        PageResponse.PaginationMetadata listPagination = PageResponse.PaginationMetadata.builder()
                .current_page(1)
                .per_page(10)
                .total_pages(1)
                .total_items(1)
                .next_page(null)
                .prev_page(null)
                .build();
                
        postListResponse = PageResponse.<List<PostResponse>>builder()
                .data(List.of(postResponse))
                .pagination(listPagination)
                .build();
    }

    @Test
    void whenCreatePost_thenReturn201AndPost() throws Exception {
        // Use the userId for the second parameter, matching what the controller will do
        when(postService.createPost(any(PostRequest.class), eq(profileId))).thenReturn(singlePostResponse);

        mockMvc.perform(post("/api/posts")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(postRequest))
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.data.id").value(postId.toString()))
                .andExpect(jsonPath("$.data.profileId").value(profileId.toString()))
                .andExpect(jsonPath("$.data.title").value("Test Post"))
                .andExpect(jsonPath("$.pagination.current_page").value(1))
                .andExpect(jsonPath("$.pagination.total_items").value(1));
    }

    @Test
    void whenGetAllPosts_thenReturn200AndPostList() throws Exception {
        when(postService.getAllPosts(1, 10)).thenReturn(postListResponse);

        mockMvc.perform(get("/api/posts"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data", hasSize(1)))
                .andExpect(jsonPath("$.data[0].id").value(postId.toString()))
                .andExpect(jsonPath("$.pagination.current_page").value(1))
                .andExpect(jsonPath("$.pagination.per_page").value(10))
                .andExpect(jsonPath("$.pagination.total_pages").value(1))
                .andExpect(jsonPath("$.pagination.total_items").value(1));
    }

    @Test
    void whenGetPostById_thenReturn200AndPost() throws Exception {
        when(postService.getPostById(postId)).thenReturn(singlePostResponse);

        mockMvc.perform(get("/api/posts/{id}", postId)
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data.id").value(postId.toString()))
                .andExpect(jsonPath("$.data.title").value("Test Post"))
                .andExpect(jsonPath("$.pagination.current_page").value(1))
                .andExpect(jsonPath("$.pagination.per_page").value(1))
                .andExpect(jsonPath("$.pagination.total_items").value(1));
    }
    
    @Test
    void whenGetMyPosts_thenReturn200AndPostList() throws Exception {
        when(postService.getPostsByProfileId(profileId, 1, 10))
                .thenReturn(postListResponse);

        mockMvc.perform(get("/api/posts/my-posts")
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data", hasSize(1)))
                .andExpect(jsonPath("$.data[0].id").value(postId.toString()))
                .andExpect(jsonPath("$.pagination.current_page").value(1))
                .andExpect(jsonPath("$.pagination.per_page").value(10));
    }
    
    @Test
    void whenUpdatePost_thenReturn200AndUpdatedPost() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(singlePostResponse);
        // Mock the updatePost call
        when(postService.updatePost(eq(postId), any(PostRequest.class))).thenReturn(singlePostResponse);

        mockMvc.perform(put("/api/posts/{id}", postId)
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(postRequest))
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data.id").value(postId.toString()));
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
                
        PageResponse<PostResponse> unauthorizedPageResponse = PageResponse.<PostResponse>builder()
                .data(unauthorizedResponse)
                .pagination(PageResponse.PaginationMetadata.builder()
                    .current_page(1)
                    .per_page(1)
                    .total_pages(1)
                    .total_items(1)
                    .build())
                .build();
                
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(unauthorizedPageResponse);

        mockMvc.perform(put("/api/posts/{id}", postId)
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(postRequest))
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isForbidden());
    }
    
    @Test
    void whenPublishPost_thenReturn200AndPublishedPost() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(singlePostResponse);
        // Mock the publishPost call
        when(postService.publishPost(postId)).thenReturn(singlePostResponse);

        mockMvc.perform(put("/api/posts/{id}/publish", postId)
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data.id").value(postId.toString()));
    }
    
    @Test
    void whenUnpublishPost_thenReturn200AndUnpublishedPost() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(singlePostResponse);
        // Mock the unpublishPost call
        when(postService.unpublishPost(postId)).thenReturn(singlePostResponse);

        mockMvc.perform(put("/api/posts/{id}/unpublish", postId)
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data.id").value(postId.toString()));
    }
    
    @Test
    void whenDeletePost_thenReturn204() throws Exception {
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(singlePostResponse);
        // Mock the deletePost call
        doNothing().when(postService).deletePost(postId);

        mockMvc.perform(delete("/api/posts/{id}", postId)
                .header("Authorization", TEST_JWT_TOKEN))
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
                
        PageResponse<PostResponse> unauthorizedPageResponse = PageResponse.<PostResponse>builder()
                .data(unauthorizedResponse)
                .pagination(PageResponse.PaginationMetadata.builder()
                    .current_page(1)
                    .per_page(1)
                    .total_pages(1)
                    .total_items(1)
                    .build())
                .build();
                
        // Mock the getPostById call for ownership check
        when(postService.getPostById(postId)).thenReturn(unauthorizedPageResponse);

        mockMvc.perform(delete("/api/posts/{id}", postId)
                .header("Authorization", TEST_JWT_TOKEN))
                .andExpect(status().isForbidden());
    }
} 