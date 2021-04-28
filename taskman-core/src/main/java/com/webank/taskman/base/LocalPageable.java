package com.webank.taskman.base;

public class LocalPageable {
    private Long startIndex;
    private Long pageSize = 10000L;

    public LocalPageable() {
    }

    public LocalPageable(Long startIndex, Long pageSize) {
        this.startIndex = startIndex;
        this.pageSize = pageSize;
    }

    public Long getPageSize() {
        return pageSize;
    }
    public void setPageSize(Long pageSize) {
        this.pageSize = pageSize;
    }
    public Long getStartIndex() {
        return startIndex;
    }
    public void setStartIndex(Long startIndex) {
        this.startIndex = startIndex;
    }
}
