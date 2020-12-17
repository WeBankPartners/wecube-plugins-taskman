package com.webank.taskman;

import org.junit.After;
import org.junit.Before;
import org.junit.runner.RunWith;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;
import org.springframework.test.context.web.WebAppConfiguration;

@RunWith(SpringRunner.class)
@SpringBootTest
@WebAppConfiguration
public class TmallApplicationTests {

    private static final Logger log = LoggerFactory.getLogger(TmallApplicationTests.class);
    @Before
    public void init() {
        log.info("start Test-----------------");
    }

    @After
    public void after() {
        log.info("end Test-----------------");
    }
}
