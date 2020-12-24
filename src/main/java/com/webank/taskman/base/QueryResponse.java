package com.webank.taskman.base;

import com.baomidou.mybatisplus.core.metadata.IPage;
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

    public QueryResponse(Long totalRows, Long startIndex, Long pageSize, List<T> contents) {
        this.pageInfo = new PageInfo(totalRows,startIndex,pageSize);
        this.contents = contents;
    }
    public QueryResponse(PageInfo pageInfo, List<T> contents) {
        this.pageInfo = pageInfo;
        this.contents = contents;
    }

    public QueryResponse(IPage<T> iPage, List<T> contents) {
        this.pageInfo = new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize());
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
