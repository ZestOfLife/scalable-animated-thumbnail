package scalable.p2.thumbnail.api;

import lombok.Getter;
import lombok.Setter;

import java.util.List;

@Setter
@Getter
public class Pojo {
    private Integer BucketID;
    private List<Video> Videos;
}

@Setter
@Getter
class Video {
    private String VideoName;
    private Integer Duration;
    private Integer FPS;
}