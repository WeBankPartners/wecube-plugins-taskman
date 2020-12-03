package com.webank.taskman.dto;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonInclude.Include;

import java.util.LinkedList;
import java.util.List;

@JsonInclude(Include.NON_NULL)
public class QueryResponse<T> {
    private PageInfo pageInfo = new PageInfo();
    private List<T> contents = new LinkedList<>();

    public QueryResponse() {
    }

    public QueryResponse(PageInfo pageInfo, List<T> contents) {
        this.pageInfo = pageInfo;
        this.contents = contents;
    }

    public PageInfo getPageInfo() {
        return pageInfo;
    }

    public void setPageInfo(PageInfo pageInfo) {
        this.pageInfo = pageInfo;
    }

    public List<T> getContents() {
        return contents;
    }

    public void setContents(List<T> domainObjs) {
        this.contents = domainObjs;
    }

    public void addContent(T ciObj) {
        this.contents.add(ciObj);
    }
}
