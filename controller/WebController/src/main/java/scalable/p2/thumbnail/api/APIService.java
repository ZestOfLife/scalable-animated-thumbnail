package scalable.p2.thumbnail.api;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.HttpClientBuilder;

import com.google.gson.Gson;
import scalable.p2.thumbnail.exceptions.InvalidCommandRequest;
import scalable.p2.thumbnail.repository.Status;
import scalable.p2.thumbnail.repository.StatusRepository;
import scalable.p2.thumbnail.statuses.Statuses;

import javax.naming.Name;
import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class APIService {
    private final StatusRepository repository;

    @Autowired
    public APIService(StatusRepository repository) {
        this.repository = repository;
    }

    private StringEntity buildRequest(JobCertificate job) throws UnsupportedEncodingException {
        Gson gson = new Gson();
        return new StringEntity(gson.toJson(job));
    }

    public void queue(Integer BucketID, String VideoName, Integer Duration, Integer FPS) throws Exception {

        HttpClient httpClient = HttpClientBuilder.create().build();
        HttpPost post = new HttpPost("event-store:8080");

        Integer DurationAt = Duration * 2/3;
        Integer length = DurationAt + 15 >= Duration ? DurationAt - Duration : 15;
        Integer ExpectedFrames = length * FPS;

        JobCertificate job = JobCertificate.builder().
                BucketID(BucketID).
                VideoName(VideoName).
                DurationAt(DurationAt * 1000).
                FPS(FPS).
                ExpectedFrames(ExpectedFrames).
                build();
        post.setEntity(buildRequest(job));
        post.setHeader("Content-type", "application/json");
        HttpResponse response = httpClient.execute(post);
        if (response.getStatusLine().getStatusCode() != 200) throw new InvalidCommandRequest(response.toString());
    }

    public Map<String, Integer> getStatus(Integer BucketID, String VideoName, Statuses state) {
        Status entry = repository.getByBucketIDAndVideoName(BucketID, VideoName);
        Map<String, Integer> ret = new HashMap<>();
        switch (state) {
            case EXTRACT -> {
                ret.put("done", entry.getExtracted());
                ret.put("percent", entry.getExtracted()/entry.getExpectedFrames() * 100);
            }
            case RESIZE -> {
                ret.put("done", entry.getResized());
                ret.put("percent", entry.getResized()/entry.getExpectedFrames() * 100);
            }
            case COMPILE -> {
                ret.put("done", entry.getCompiled());
                ret.put("percent", entry.getCompiled()/entry.getExpectedFrames() * 100);
            }
        }
        return ret;
    }

    public List<Status> getGifs(Integer BucketID) {
        return repository.getCompletedGifs(BucketID);
    }

}
