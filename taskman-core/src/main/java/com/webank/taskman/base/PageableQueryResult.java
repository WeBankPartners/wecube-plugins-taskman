package com.webank.taskman.base;

import java.util.LinkedList;
import java.util.List;

import com.baomidou.mybatisplus.core.metadata.IPage;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonInclude.Include;

@JsonInclude(Include.NON_NULL)
public class PageableQueryResult<T> {
    private LocalPageInfo pageInfo = new LocalPageInfo();
    private List<T> contents = new LinkedList<>();

    public PageableQueryResult() {
    }

    public PageableQueryResult(Long totalRows, Long startIndex, Long pageSize, List<T> contents) {
        this.pageInfo = new LocalPageInfo(totalRows,startIndex,pageSize);
        this.contents = contents;
    }
    public PageableQueryResult(LocalPageInfo pageInfo, List<T> contents) {
        this.pageInfo = pageInfo;
        this.contents = contents;
    }

    public PageableQueryResult(IPage<T> iPage, List<T> contents) {
        this.pageInfo = new LocalPageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize());
        this.contents = contents;
    }


    public LocalPageInfo getPageInfo() {
        return pageInfo;
    }

    public void setPageInfo(LocalPageInfo pageInfo) {
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
