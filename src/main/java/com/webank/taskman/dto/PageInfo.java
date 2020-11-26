package com.webank.taskman.dto;

public class PageInfo extends Pageable {
	private int totalRows;

    public PageInfo() {
        this.totalRows = 0;
    }

	public PageInfo(int totalRows, int startIndex, int pageSize) {
		super(startIndex, pageSize);
		this.totalRows = totalRows;
	}

	public int getTotalRows() {
		return totalRows;
	}

	public void setTotalRows(int totalRows) {
		this.totalRows = totalRows;
	}
}
