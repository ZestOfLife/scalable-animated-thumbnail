package scalable.p2.thumbnail.api;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Builder
public class JobCertificate {
    private Integer BucketID;
    private String VideoName;
    private Integer ExpectedFrames;
    private Integer FPS;
    private Integer DurationAt;
}
