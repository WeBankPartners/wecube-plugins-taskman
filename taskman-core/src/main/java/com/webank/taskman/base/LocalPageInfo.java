package com.webank.taskman.base;


public class LocalPageInfo extends LocalPageable {
    private Long totalRows;

    public LocalPageInfo() {
        this.totalRows = 0L;
    }

    public LocalPageInfo(Long totalRows, Long startIndex, Long pageSize) {
        super(startIndex, pageSize);
        this.totalRows = totalRows;
    }

    public Long getTotalRows() {
        return totalRows;
    }

    public void setTotalRows(Long totalRows) {
        this.totalRows = totalRows;
    }
}
