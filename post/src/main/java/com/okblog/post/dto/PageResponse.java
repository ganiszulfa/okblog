package com.okblog.post.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class PageResponse<T> {
    
    private T data;
    private PaginationMetadata pagination;
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class PaginationMetadata {
        private int current_page;
        private int per_page;
        private int total_pages;
        private long total_items;
        private Integer next_page;
        private Integer prev_page;
    }
} 