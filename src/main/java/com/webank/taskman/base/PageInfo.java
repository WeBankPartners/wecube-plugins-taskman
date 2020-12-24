package com.webank.taskman.base;

public class PageInfo extends Pageable {
	private Long totalRows;

    public PageInfo() {
        this.totalRows = 0L;
    }

	public PageInfo(Long totalRows, Long startIndex, Long pageSize) {
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
